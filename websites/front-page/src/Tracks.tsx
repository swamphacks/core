import "./Tracks.css";
import controller from "./assets/Controller.png";
import Gator from "./assets/gator.gif";
import { useState } from "react";
import Modal from "react-modal";
import Earth from "./assets/earth.png";
import Robot from "./assets/robot.png";
import Computer from "./assets/computer.png";
import Apple from "./assets/apple.png";

const tracks = [
  {
    name: "Education, Accessibility & Social Impact",
    icon: Earth,
    description:
      "Technology for the people. Create solutions that improve accessibility, empower communities, expand educational opportunities, or tackle real-world challenges.",
  },
  {
    name: "Artificial Intelligence",
    icon: Robot,
    description:
      "Push the boundaries of what's possible. Leverage AI to build projects powered by artificial intelligence, from LLMs and AI agents to computer vision and machine learning that make an impact across any domain.",
  },
  {
    name: "Hardware & Physical Computing (Get Physical With It)",
    icon: Computer,
    description:
      "Bring ideas into the physical world. Build hardware-based projects or develop software that interacts, controls, or enhances physical devices and embedded systems.",
  },
  {
    name: "Entertainment & Creative Technology",
    icon: controller,
    description:
      "Let creativity take the lead. Create engaging digital experiences through games, interactive media, digital art, storytelling, music, and showcase the power of technology and imagination.",
  },
  {
    name: "Health & Wellness",
    icon: Apple,
    description:
      "Support healthier minds and healthier lives. Innovate technology that promotes mental health and physical well-being by supporting emotional wellness and healthy habits. General health-focused projects are also welcome.",
  },
];

const customModalStyles = {
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
    zIndex: "100",
  },
  content: {
    overflow: "hidden",
    borderColor: "#231909",
    borderWidth: "5px",
    backgroundColor: "#7D573C",
    top: "50%",
    left: "50%",
    right: "auto",
    bottom: "auto",
    marginRight: "-50%",
    transform: "translate(-50%, -50%)",
    zIndex: "100",
  },
};

export default function Tracks() {
  const [selectedTrack, setSelectedTrack] = useState<
    (typeof tracks)[number] | null
  >(null);

  function closeModal() {
    setSelectedTrack(null);
  }

  return (
    <div className="tracks-container">
      <h1 className="tracks-header">Tracks</h1>

      <div className="tracks-body-container">
        <div className="track-container">
          {tracks.map((track) => (
            <div
              key={track.name}
              className="track"
              onClick={() => setSelectedTrack(track)}
            >
              <img className="track-icon" src={track.icon} />
              <span className="track-header">{track.name}</span>
            </div>
          ))}
        </div>

        <img className="track-gator" src={Gator} />
      </div>

      <Modal
        isOpen={Boolean(selectedTrack)}
        onRequestClose={closeModal}
        style={customModalStyles}
      >
        <div className="tracks-modal-container">
          <button
            onClick={closeModal}
            className="modal-close-btn nes-btn is-error"
          >
            X
          </button>
          <div>
            <div className="modal-track">
              <img className="track-icon" src={selectedTrack?.icon} />
              <p className="track-title">{selectedTrack?.name}</p>
            </div>

            <p className="modal-track-description">
              {selectedTrack?.description}
            </p>
          </div>
        </div>
      </Modal>
    </div>
  );
}
