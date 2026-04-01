import { Rating } from "@smastrom/react-rating";

interface RatingFieldsProps {
  experience: number;
  passion: number;
  onExperience: (value: number) => void;
  onPassion: (value: number) => void;
}

export const RatingFields = ({
  experience,
  passion,
  onExperience,
  onPassion,
}: RatingFieldsProps) => {
  return (
    <div className="w-1/3">
      <div className="flex flex-row gap-4 justify-between items-center">
        <p className="text-lg">Experience:</p>
        <Rating
          style={{ maxWidth: 150 }}
          value={experience}
          onChange={onExperience}
        />
      </div>

      <div className="flex flex-row gap-4 items-center justify-between mt-4">
        <p className="text-lg">Passion:</p>
        <Rating
          style={{ maxWidth: 150 }}
          value={passion}
          onChange={onPassion}
        />
      </div>
    </div>
  );
};
