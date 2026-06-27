type PageHeadingProps = {
  title: string;
  description: string;
};

export function PageHeading({ title, description }: PageHeadingProps) {
  return (
    <div className="space-y-2">
      <h1 className="text-2xl font-semibold text-white">{title}</h1>
      <p className="text-sm text-[#CBD5E1]">{description}</p>
    </div>
  );
}
