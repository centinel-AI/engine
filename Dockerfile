ARG ALPINE_VERSION=3.21.3

FROM alpine:${ALPINE_VERSION}

ARG TARGETARCH
# terraform | opentofu
ARG ENGINE
ARG TERRAFORM_VERSION
ARG OPENTOFU_VERSION
# azure | gcp | aws | oci — directory under providers/ copied into the image
ARG CLOUD_PROVIDER
ENV ENGINE_PROVIDER=${CLOUD_PROVIDER}

ENV TF_PLUGIN_CACHE_DIR="/app/workspace/.terraform.d/plugin-cache"

RUN apk add --no-cache \
    curl \
    unzip \
    bash

COPY scripts/inject-provider-versions.sh /tmp/inject-provider-versions.sh
RUN chmod +x /tmp/inject-provider-versions.sh

RUN if [ "$ENGINE" = "opentofu" ]; then \
      if [ -z "$OPENTOFU_VERSION" ]; then echo "OPENTOFU_VERSION is required when ENGINE=opentofu" >&2; exit 1; fi; \
      OPENTOFU_VERSION_TRIM="${OPENTOFU_VERSION#v}"; \
      curl -fsSL "https://github.com/opentofu/opentofu/releases/download/v${OPENTOFU_VERSION_TRIM}/tofu_${OPENTOFU_VERSION_TRIM}_linux_${TARGETARCH}.zip" \
        -o engine.zip && \
      unzip engine.zip && \
      mv tofu /usr/local/bin/tofu && \
      rm engine.zip && \
      chmod +x /usr/local/bin/tofu && \
      echo -n tofu > /tmp/iac_engine_label; \
    else \
      if [ -z "$TERRAFORM_VERSION" ]; then echo "TERRAFORM_VERSION is required when ENGINE=terraform" >&2; exit 1; fi; \
      curl -fsSL "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip" \
        -o engine.zip && \
      unzip engine.zip && \
      mv terraform /usr/local/bin/terraform && \
      rm engine.zip && \
      chmod +x /usr/local/bin/terraform && \
      echo -n terraform > /tmp/iac_engine_label; \
    fi

COPY entrypoint.sh /app/workspace/entrypoint.sh
RUN chmod +x /app/workspace/entrypoint.sh

COPY providers/${CLOUD_PROVIDER}/ /app/workspace/

COPY scripts/sync-backend-config.sh /app/workspace/sync-backend-config.sh
RUN chmod +x /app/workspace/sync-backend-config.sh

WORKDIR /app/workspace

RUN mkdir -p /app/workspace/.terraform.d/plugin-cache \
    && mkdir -p /app/workspace/data \
    && mv /tmp/iac_engine_label /app/workspace/.iac_engine_bin

ARG AZURERM_VERSION
ARG AZUREAD_VERSION
ARG GOOGLE_PROVIDER_VERSION
ARG AWS_PROVIDER_VERSION
ARG OCI_PROVIDER_VERSION
ARG TERRAFORM_REQUIRED_VERSION
ARG OPENTOFU_REQUIRED_VERSION

RUN /tmp/inject-provider-versions.sh /app/workspace

RUN ENGINE_BIN="$(tr -d '\n\r' < /app/workspace/.iac_engine_bin)" && "$ENGINE_BIN" init

ENTRYPOINT ["/app/workspace/entrypoint.sh"]
CMD ["terraform", "plan"]
