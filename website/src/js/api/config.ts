import type { JSONSchemaType, DefinedError, ValidateFunction } from "ajv";
import Ajv from "ajv";
import axios from "axios";

export interface Config {
  oidc_server_url: string;
  oidc_server_public_client_id: string;
  oidc_server_discovery_endpoint: string;
}

const ConfigSchema: JSONSchemaType<Config> = {
  type: "object",
  properties: {
    oidc_server_url: { type: "string", nullable: false },
    oidc_server_public_client_id: { type: "string", nullable: false },
    oidc_server_discovery_endpoint: { type: "string", nullable: false },
  },
  required: [
    "oidc_server_url",
    "oidc_server_public_client_id",
    "oidc_server_discovery_endpoint",
  ],
};

const ajv = new Ajv();

export async function getConfig(): Promise<Config> {
  const response = await axios.get<Config>("/v1/public/spa-config");
  if (response.status !== 200) {
    throw new Error(`get config failed with ${response.status} status code`);
  }
  const validateConfig = ajv.compile(ConfigSchema);
  if (validateConfig(response.data)) {
    return response.data;
  }
  throw new Error(getAjvErrorMessage(validateConfig));
}

function getAjvErrorMessage(validate: ValidateFunction): string {
  try {
    let errorMessage = "";
    for (const err of validate.errors as DefinedError[]) {
      switch (err.keyword) {
        case "type":
          errorMessage += err.message + ", ";
          break;
      }
    }
    return `Invalidate response type: ${errorMessage}`;
  } catch (err) {
    return "Invalidate response type";
  }
}
