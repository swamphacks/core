import "./Tracks.css";
import pencil from "./assets/Pencil.png";
import controller from "./assets/Controller.png";
import star from "./assets/Star.png";
import leaf from "./assets/Leaf.png";
import handshake from "./assets/Handshake.png";

const tracks = [
  {
    name: "Overall Prize",
    icon: star,
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "Education, Accessibility & Social Impact",
    icon: pencil,
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "Sustainability",
    icon: leaf,
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "Game Design",
    icon: controller,
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
  {
    name: "The Design of Everyday Life (Human-Centered Design)",
    icon: handshake,
    description:
      "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
  },
];

export default function Tracks() {
  // const [openIndex, setOpenIndex] = useState<number | null>(null);

  // const toggle = (index: number) => {
  //   setOpenIndex(openIndex === index ? null : index);
  // };

  return (
    <div className="tracks-container">
      <h1 className="tracks-header">Tracks</h1>

      <div className="track-container">
        {tracks.map((track) => (
          <div
            key={track.name}
            className="track"
            // onClick={() => toggle(idx)}
          >
            <img className="track-icon" src={track.icon} />
            <span className="track-header">{track.name}</span>

            {/* {openIndex === idx && (
              <div className="track-description">
                <p>{track.description}</p>
              </div>
            )} */}
          </div>
        ))}
      </div>
      {/* <div className="track-items-container">
        <div className="track-icons-container">
          {tracks.map((track) => (
            <img className="track-icon" src={track.icon} />
          ))}
        </div>

        <div className="track-names-container">
          {tracks.map((track) => (
            <span className="track-header">{track.name}</span>
          ))}
        </div>
      </div> */}
    </div>
  );
}
