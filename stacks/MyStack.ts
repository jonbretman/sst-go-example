import { RDS, Api, StackContext } from "@serverless-stack/resources";

export function MyStack({ stack }: StackContext) {
  // Create the Aurora DB cluster
  const cluster = new RDS(stack, "Cluster", {
    engine: "postgresql10.14",
    defaultDatabaseName: "postgres",
    migrations: "migrations",
  });

  // Create the HTTP API
  const api = new Api(stack, "Api", {
    defaults: {
      function: {
        environment: {
          DATABASE_NAME: "postgres",
          DATABASE_RESOURCE_ARN: cluster.clusterArn,
          DATABASE_SECRET_ARN: cluster.secretArn,
        },
        permissions: [cluster],
      },
    },
    routes: {
      "GET /": "services/hello",
    },
  });

  // Show API endpoint in output
  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}
