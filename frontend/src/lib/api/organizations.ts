import { apiClient } from "@/lib/api/api-client";
import { apiRoutes } from "@/lib/api/api-routes";
import type { Organization } from "@/lib/api/types";

export async function listOrganizations(): Promise<Organization[]> {
  if (!apiClient.defaults.baseURL) {
    return [
      {
        id: "1",
        name: "Acme Corp",
        slug: "acme",
        createdAt: new Date().toISOString(),
      },
      {
        id: "2",
        name: "Globex",
        slug: "globex",
        createdAt: new Date().toISOString(),
      },
    ];
  }

  const { data } = await apiClient.get<Organization[]>(apiRoutes.organizations);
  return data;
}
