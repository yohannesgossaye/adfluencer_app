import { InfluencerProfile } from "@/features/business/influencer-profile";

export default async function BusinessInfluencerProfilePage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return <InfluencerProfile influencerId={id} />;
}
