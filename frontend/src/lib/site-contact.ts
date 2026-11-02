/**
 * How visitors reach the Adfluence team. Change these three values when the
 * company email, phone line and Telegram account are ready; every page
 * (home "Talk to us", footer, pricing) reads them from here.
 * The labels shown on the site stay neutral ("Email us", "Call us").
 */
export const SITE_CONTACT = {
  email: "yohannesgossaye73@gmail.com",
  phone: "+251953422606",
  telegram: "jogomez11",
} as const;

export const contactLinks = {
  email: `mailto:${SITE_CONTACT.email}`,
  phone: `tel:${SITE_CONTACT.phone}`,
  telegram: `https://t.me/${SITE_CONTACT.telegram}`,
};
