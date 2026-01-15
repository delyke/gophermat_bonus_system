FROM --platform=linux/amd64 debian:bookworm-slim

RUN useradd -m appuser
WORKDIR /app
USER appuser

COPY cmd/accrual/accrual_linux_amd64 ./accrual

EXPOSE 8080
ENTRYPOINT ["./accrual"]