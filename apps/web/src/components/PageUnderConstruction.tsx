import TablerTools from "~icons/tabler/tools";

export function PageUnderConstruction() {
  return (
    <div className="flex w-full h-full items-center justify-center">
      <div className="flex items-center gap-2">
        <TablerTools className="text-xl" />
        <h1 className="text-xl font-medium">Page Under Construction</h1>
      </div>
    </div>
  );
}
