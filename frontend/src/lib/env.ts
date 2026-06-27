import { z } from "zod";

const publicEnvSchema = z.object({
  NEXT_PUBLIC_API_URL: z.string().optional(),
});

export type PublicEnv = z.infer<typeof publicEnvSchema> & {
  apiBaseUrl: string;
};

export function getPublicEnv(): PublicEnv {
  const parsed = publicEnvSchema.safeParse({
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
  });

  const apiUrl = parsed.success ? parsed.data.NEXT_PUBLIC_API_URL : undefined;

  return {
    NEXT_PUBLIC_API_URL: apiUrl,
    apiBaseUrl: apiUrl?.replace(/\/$/, "") ?? "",
  };
}
