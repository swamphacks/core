import { useState } from "react";
import { useDebounce } from "@uidotdev/usehooks";
import { Input } from "@/components/ui/Field";
import ApplicationList from "@/modules/Application/ApplicationSearch/ApplicationList";

export default function ApplicationSearchPage() {
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearchInput = useDebounce(searchInput, 500);

  return (
    <div>
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Application search</h1>
          <p className="text-sm text-muted-foreground">
            Browse applications and inspect status for each registered user.
          </p>
        </div>
      </div>

      <div className="mt-5 p-2 rounded-md bg-surface max-w-fit flex flex-col">
        <div className="w-70 mb-3">
          <Input
            aria-label="search input"
            className="rounded-md"
            placeholder="Search by username or email"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
          />
        </div>
        <ApplicationList searchInput={debouncedSearchInput} />
      </div>
    </div>
  );
}
