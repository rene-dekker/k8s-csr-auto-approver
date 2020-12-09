FROM scratch

COPY init-container /init-container
ENTRYPOINT ["/init-container"]
