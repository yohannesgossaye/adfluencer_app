import { CreateCampaign } from "@/features/business/create-campaign";

export default async function BusinessEditCampaignPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <CreateCampaign campaignId={id} />;
}
