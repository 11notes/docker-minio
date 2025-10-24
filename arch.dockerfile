# ╔═════════════════════════════════════════════════════╗
# ║                       SETUP                         ║
# ╚═════════════════════════════════════════════════════╝
# GLOBAL
  ARG APP_UID=1000 \
      APP_GID=1000 \
      BUILD_ROOT=/go/minio \
      BUILD_SRC=minio/minio.git \
      BUILD_BIN=/minio

# :: FOREIGN IMAGES
  FROM 11notes/distroless AS distroless
  FROM 11notes/distroless:localhealth AS distroless-localhealth
  FROM 11notes/distroless:openssl AS distroless-openssl

# ╔═════════════════════════════════════════════════════╗
# ║                       BUILD                         ║
# ╚═════════════════════════════════════════════════════╝
# :: MINIO
  FROM 11notes/go:1.25 AS build
  ARG APP_VERSION \
      APP_VERSION_BUILD \
      BUILD_ROOT \
      BUILD_SRC \
      BUILD_BIN

  RUN set -ex; \
    SEMVER=$(echo ${APP_VERSION} | sed 's|\.|-|g'); \
    eleven git clone ${BUILD_SRC} RELEASE.${SEMVER}T${APP_VERSION_BUILD}; \
    sed -i 's|"DEVELOPMENT.GOGET"|"RELEASE.'${SEMVER}'T'${APP_VERSION_BUILD}'"|g' ${BUILD_ROOT}/cmd/build-constants.go; \
    sed -i 's|fmt.Fprintln(banner, color.Blue("Version:")+color.Bold(" %s (%s %s/%s)", ReleaseTag, runtime.Version(), runtime.GOOS, runtime.GOARCH))|fmt.Fprint(banner, color.Blue("Version:")+color.Bold(" %s", ReleaseTag))|g' ${BUILD_ROOT}/cmd/main.go; \
    sed -i 's|fmt.Fprintln(banner, color.Blue("Copyright:|//|g' ${BUILD_ROOT}/cmd/main.go; \
    sed -i 's|fmt.Fprintln(banner, color.Blue("License:|//|g' ${BUILD_ROOT}/cmd/main.go; \
    sed -i 's|logger.Startup(color.Blue("\\nDocs:|//|g' ${BUILD_ROOT}/cmd/server-startup-msg.go;

  RUN set -ex; \
    cd ${BUILD_ROOT}; \
    eleven go build ${BUILD_BIN} main.go;

  RUN set -ex; \
    eleven distroless ${BUILD_BIN};

# :: ENTRYPOINT
  FROM 11notes/go:1.25 AS entrypoint
  COPY ./build /
  ARG APP_VERSION \
      APP_VERSION_BUILD

  RUN set -ex; \
    cd /go/entrypoint; \
    eleven go build /entrypoint main.go;

  RUN set -ex; \
    eleven distroless /entrypoint;

# :: FILE SYSTEM
  FROM alpine AS file-system
  ARG APP_ROOT \
      APP_UID \
      APP_GID

  RUN set -ex; \
    mkdir -p /distroless${APP_ROOT}/ssl/CAs; \
    mkdir -p /distroless/mnt;

# ╔═════════════════════════════════════════════════════╗
# ║                       IMAGE                         ║
# ╚═════════════════════════════════════════════════════╝
# :: HEADER
  FROM scratch

  # :: default arguments
    ARG TARGETPLATFORM \
        TARGETOS \
        TARGETARCH \
        TARGETVARIANT \
        APP_IMAGE \
        APP_NAME \
        APP_VERSION \
        APP_ROOT \
        APP_UID \
        APP_GID \
        APP_NO_CACHE

  # :: default environment
    ENV APP_IMAGE=${APP_IMAGE} \
        APP_NAME=${APP_NAME} \
        APP_VERSION=${APP_VERSION} \
        APP_ROOT=${APP_ROOT}

  # :: app specific environment
    ENV MINIO_ROOT_USER="admin" \
        MINIO_BROWSER_LOGIN_ANIMATION="off"

  # :: multi-stage
    COPY --from=distroless / /
    COPY --from=distroless-localhealth / /
    COPY --from=distroless-openssl / /
    COPY --from=build /distroless/ /
    COPY --from=entrypoint /distroless/ /
    COPY --from=file-system --chown=${APP_UID}:${APP_GID} /distroless/ /

# :: PERSISTENT DATA
  VOLUME ["${APP_ROOT}/ssl"]

# :: MONITORING
  HEALTHCHECK --interval=5s --timeout=2s --start-period=5s \
    CMD ["/usr/local/bin/localhealth", "https://127.0.0.1:9001/"]

# :: EXECUTE
  USER ${APP_UID}:${APP_GID}
  ENTRYPOINT ["/usr/local/bin/entrypoint"]