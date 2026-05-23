interface NavSectionProps {
  name: string;
}

export function NavSection({ name }: NavSectionProps) {
  return (
    <div className="mt-5 py-1 flex items-center pl-3.5">
      <p className="text-text-secondary text-sm opacity-70">{name}</p>
      <hr className="h-[1px] w-full bg-navlink-bg-active border-0 ml-2"></hr>
    </div>
  );
}
