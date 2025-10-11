import TablerLoader from "~icons/tabler/loader-2";

export function PageLoading() {
  return (
    <div className="flex w-full h-full items-center justify-center">
      <div className="flex items-center gap-2">
        <TablerLoader className="text-xl animate-spin" />
        <h1 className="text-xl font-medium">Loading</h1>
      </div>
    </div>
  );
}
