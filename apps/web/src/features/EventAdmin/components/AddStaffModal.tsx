import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { TextField } from "@/components/ui/TextField";
import { useUsers } from "@/features/Users/hooks/useUsers";
import { useState } from "react";

const AddStaffModal = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const { data, isLoading, isError } = useUsers(searchTerm)

  // Tanstack Query can be used here to fetch users based on the search term

  return (
    <Modal size="2xl" title="Add Staff" className="h-1/2">
      <div className="w-full h-full py-4 flex flex-col gap-6">
        <TextField
          className="w-full"
          placeholder="Search by name or email..."
          onChange={(e) => setSearchTerm(e)}
        />

        {isLoading && <p>Loading...</p>}
        {isError && <p>Error loading users.</p>}
        {data && data.length === 0 && <p>No users found.</p>}
        {data && data.length > 0 && (
          <div className="flex flex-col gap-4 max-h-96 overflow-y-auto">
            {data.map((user) => (
              <div
                key={user.id}
                className="flex flex-row justify-between items-center p-4 border border-border rounded-md"
              >
                <div className="flex flex-col">
                  <p className="font-medium">{user.name}</p>
                  <p className="text-sm text-muted-foreground">{user.email}</p>
                </div>
                <Button variant="primary">Add</Button>
              </div>
            ))}
          </div>
        )}

        <p>Searching for {searchTerm}</p>
      </div>
    </Modal>
  );
};

export default AddStaffModal;
