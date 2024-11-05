FROM denoland/deno

WORKDIR /app

COPY server.ts key_value_routes.ts view_routes.ts replicate.ts internal_routes.ts ./

RUN deno install

# App is already listening on port 8090
CMD ["deno", "run", "--allow-net", "server.ts"]

EXPOSE 8090:8090