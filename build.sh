#! /bin/sh

set -e

PLATFORM="linux/amd64,linux/arm64"
REPOSITORY="jlmbrt/debug-server"
TAG=${TAG:-$(git rev-parse --short HEAD)}

PUSH=0
LATEST=0

IMAGE="${REPOSITORY}:${TAG}"
IMAGE_LATEST="${REPOSITORY}:latest"

for arg in "$@"; do
	case "$arg" in
	--push) PUSH=1 ;;
	--latest) LATEST=1 ;;
	esac
done

docker build \
	--platform "${PLATFORM}" \
	-t "${IMAGE}" \
	--output "type=docker" \
	.

[ "$PUSH" -eq 1 ] && docker push "${IMAGE}"

[ "$PUSH" -eq 1 ] && [ "$LATEST" -eq 1 ] && docker tag "${IMAGE}" "${IMAGE_LATEST}" && docker push "${IMAGE_LATEST}"

echo "${IMAGE}"
