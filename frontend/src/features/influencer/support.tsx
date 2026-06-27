"use client";

import {
  influencerIssueTypeOptions,
  relatedInfluencerCampaignOptions,
} from "@/data/support";
import { UserSupportPage } from "@/features/support/user-support-page";

export function InfluencerSupport() {
  return (
    <UserSupportPage
      role="influencer"
      subtitle="Contact Adfluence support for profile, campaign, contract, or payout issues."
      issueTypeOptions={influencerIssueTypeOptions}
      relatedOptions={relatedInfluencerCampaignOptions}
    />
  );
}
