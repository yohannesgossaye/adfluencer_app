"use client";

import {
  businessIssueTypeOptions,
  relatedBusinessCampaignOptions,
} from "@/data/support";
import { UserSupportPage } from "@/features/support/user-support-page";

export function BusinessSupport() {
  return (
    <UserSupportPage
      role="business"
      subtitle="Get help with campaigns, contracts, payments, or creator collaborations."
      issueTypeOptions={businessIssueTypeOptions}
      relatedOptions={relatedBusinessCampaignOptions}
    />
  );
}
