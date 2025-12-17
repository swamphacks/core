import { Button } from "@/components/ui/Button";
import WeightField from "./WeightField";
import { useState } from "react";

interface ApplicationReviewPageProps {
  eventId: string;
}

export default function ApplicationDecisionsPage({
  eventId,
}: ApplicationReviewPageProps) {
  const [experienceWeight, setExperienceWeight] = useState(0);
  const [passionWeight, setPassionWeight] = useState(0);

  const handleCalculate = () => {
    console.log(eventId);
  };

  return (
    <div className="pb-8">
      <div className="flex flex-row gap-2 items-center">
        <WeightField
          name="Experience weight:"
          value={experienceWeight}
          setValue={setExperienceWeight}
        />
        <WeightField
          name="Passion weight:"
          value={passionWeight}
          setValue={setPassionWeight}
        />
      </div>
      <div className="flex flex-row gap-2 items-center">
        <Button className="w-fit" varient="primary" onPress={handleCalculate}>
          Calculate
        </Button>
      </div>
    </div>
  );
}
