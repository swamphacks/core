import TablerLoader from "~icons/tabler/loader-2";

export default function Loading() {
  return (
    <div className="flex items-center gap-2">
      <TablerLoader className="text-xl animate-spin" />
      <h1 className="text-xl font-medium">Loading</h1>
    </div>
  );
}
