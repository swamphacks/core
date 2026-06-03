import { useState } from "react";

const tracks = [
  {
    name: "Overall Prize",
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "Education, Accessibility & Social Impact",
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "Sustainability",
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "Game Design",
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "The Design of Everyday Life (Human-Centered Design)",
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
];

export function Tracks() {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggle = (index: number) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <div>
      <p className="tab-title">Tracks</p>
      <div className="tab-container track-container">
        {tracks.map((track, idx) => (
          <div
            key={track.name}
            className={`track ${openIndex === idx ? "open" : ""}`}
          >
            <button
              className={`track-header`}
              onClick={() => toggle(idx)}
              aria-expanded={openIndex === idx}
            >
              {track.name}
            </button>

            {openIndex === idx && (
              <div className="track-description">
                <p>{track.description}</p>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
