${{ content_synopsis }} This image will run minio [rootless](https://github.com/11notes/RTFM/blob/main/linux/container/image/rootless.md) and [distroless](https://github.com/11notes/RTFM/blob/main/linux/container/image/distroless.md), for maximum security and performance. In addition to being small and secure, it will also automatically create the needed SSL certificate based on the Root CA you add via compose. Checkout [mc](https://github.com/11notes/docker-mc) to manage your cluster from CLI or use [minio-console](https://github.com/11notes/docker-minio-console) to manage it via the full web UI (pre cull).

${{ content_uvp }} Good question! Because ...

${{ github:> [!IMPORTANT] }}
${{ github:> }}* ... this image runs [rootless](https://github.com/11notes/RTFM/blob/main/linux/container/image/rootless.md) as 1000:1000
${{ github:> }}* ... this image has no shell since it is [distroless](https://github.com/11notes/RTFM/blob/main/linux/container/image/distroless.md)
${{ github:> }}* ... this image is auto updated to the latest version via CI/CD
${{ github:> }}* ... this image has a health check
${{ github:> }}* ... this image runs read-only
${{ github:> }}* ... this image is automatically scanned for CVEs before and after publishing
${{ github:> }}* ... this image is created via a secure and pinned CI/CD process
${{ github:> }}* ... this image is very small

If you value security, simplicity and optimizations to the extreme, then this image might be for you.

${{ content_comparison }}

${{ title_volumes }}
* **${{ json_root }}/ssl** - Directory of SSL certificates

${{ content_compose }}

${{ content_defaults }}
| `--address` | 0.0.0.0:9000 | minio IP and port |
| `--console-address` | 0.0.0.0:3000 | minio console IP and port |
| `--certs-dir` | ${{ json_root }}/ssl | minio SSL directory |
| `--anonymous` |  | hide sensitiv data in log output |
| `--json` |  | log as json |

${{ content_environment }}
| `MINIO_ROOT_USER` | username of admin account | admin |
| `MINIO_BROWSER_LOGIN_ANIMATION` | enable console login animation | off |

${{ content_source }}

${{ content_parent }}

${{ content_built }}

${{ content_tips }}