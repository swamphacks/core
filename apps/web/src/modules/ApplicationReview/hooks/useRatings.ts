import { useEffect, useState } from "react";

export const useRatings = (
  initialExperience: number,
  initialPassion: number,
) => {
  const [experience, setExperience] = useState(initialExperience);
  const [passion, setPassion] = useState(initialPassion);
  const [isDirty, setDirty] = useState(false);

  useEffect(() => {
    setDirty(experience !== initialExperience || passion !== initialPassion);
  }, [experience, passion, initialExperience, initialPassion]);

  useEffect(() => {
    setExperience(initialExperience);
    setPassion(initialPassion);
  }, [initialExperience, initialPassion]);

  const reset = () => {
    setExperience(initialExperience);
    setPassion(initialPassion);
  };

  return {
    experience,
    passion,
    isDirty,
    setExperience,
    setPassion,
    reset,
  };
};
