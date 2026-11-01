import { CampaignDetailPage } from "@/features/business/campaign-detail";

export default async function BusinessCampaignPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <CampaignDetailPage campaignId={id} />;
}
