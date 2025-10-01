import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { TextField } from "@/components/ui/TextField";
import { useEventStaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import { useUsers } from "@/features/Users/hooks/useUsers";
import { useCallback, useState } from "react";
import debounce from "lodash.debounce";
import { cn } from "@/utils/cn";

interface Props {
  eventId: string;
}

const AddStaffModal = ({ eventId }: Props) => {
  const [query, setQuery] = useState<string | null>("");
  const [selectedUserIds, setSelectedUserIds] = useState<string[]>([]);

  const {
    data: users,
    isLoading: isUsersLoading,
    isError: isUsersError,
  } = useUsers(query);
  const {
    data: eventStaff,
    isLoading: isStaffUsersLoading,
    isError: isStaffUsersError,
  } = useEventStaffUsers(eventId);

  // Tanstack Query can be used here to fetch users based on the search term

  const handleSearch = useCallback(
    debounce((value: string) => {
      setQuery(value);
    }, 500),
    [],
  );

  const onSelectionChange = (userId: string) => {
    setSelectedUserIds((prev) => {
      if (prev.includes(userId)) {
        return prev.filter((id) => id !== userId);
      } else {
        return [...prev, userId];
      }
    });
  };

  const isAlreadyStaff = (userId: string) =>
    eventStaff?.some((staff) => staff.id === userId);
  const isSelected = (userId: string) => selectedUserIds.includes(userId);

  return (
    <Modal size="2xl" title="Add Staff" className="h-1/2">
      <div className="w-full h-full justify-between py-4 flex flex-col gap-6">
        <div className="flex-1 flex flex-col gap-6">
          <TextField
            className="w-full"
            aria-label="Search users"
            placeholder="Search by name or email..."
            onChange={(e) => handleSearch(e)}
          />

          {query === "" ? (
            <p className="text-muted-foreground">Type to search for users</p>
          ) : isUsersLoading || isStaffUsersLoading ? (
            <p>Loading...</p>
          ) : isUsersError || isStaffUsersError ? (
            <p>Error loading users.</p>
          ) : users && users.length === 0 ? (
            <p>No users found.</p>
          ) : (
            users && (
              <div className="flex flex-col gap-4 max-h-96 overflow-y-auto">
                {users.map((user) => {
                  const alreadyStaff = isAlreadyStaff(user.id);
                  const selected = isSelected(user.id);

                  return (
                    <div
                      onClick={() =>
                        !alreadyStaff && onSelectionChange(user.id)
                      }
                      className={cn(
                        "w-full bg-input-bg hover:cursor-pointer hover:bg-neutral-900 rounded-md p-4 flex flex-row items-center justify-between select-none",
                        selected
                          ? "border-1 border-green-700"
                          : alreadyStaff
                            ? "opacity-50 cursor-not-allowed hover:cursor-default hover:bg-input-bg"
                            : "border-1 border-transparent",
                      )}
                    >
                      <div className="flex flex-col">
                        <p className="font-medium">{user.name}</p>
                        <p className="text-sm text-text-secondary">
                          {user.email}
                        </p>
                      </div>

                      {alreadyStaff && (
                        <p className="text-sm text-text-secondary font-medium">
                          Already Staff
                        </p>
                      )}
                    </div>
                  );
                })}
              </div>
            )
          )}
        </div>

        {selectedUserIds.length > 0 && (
          <Button variant="primary">
            Add {selectedUserIds.length} Staff Member
            {selectedUserIds.length > 1 ? "s" : ""}
          </Button>
        )}
      </div>
    </Modal>
  );
};

export default AddStaffModal;
