// Demo-only auth users. Replace with backend auth before production.
export type MockUserRole = "business" | "influencer" | "admin";

export type MockUserRecord = {
  email: string;
  password: string;
  role: MockUserRole;
  name: string;
  redirect: string;
};

export const mockUsers: MockUserRecord[] = [
  {
    email: "yohannes@gmail.com",
    password: "yoya1111",
    role: "business",
    name: "Yohannes",
    redirect: "/business/dashboard",
  },
  {
    email: "biruk@gmail.com",
    password: "bg1114",
    role: "influencer",
    name: "Biruk",
    redirect: "/influencer/dashboard",
  },
  {
    email: "admin@adfluence.com",
    password: "admin1111",
    role: "admin",
    name: "Admin",
    redirect: "/admin/dashboard",
  },
];

