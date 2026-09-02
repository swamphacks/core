import "./Tracks.css";
import controller from "./assets/Controller.png";
import Gator from "./assets/gator.gif";
import { useState } from "react";
import Modal from "react-modal";
import Robot from "./assets/robot.png";
import Computer from "./assets/computer.png";
import Star from "./assets/Star.png";
import Handshake from "./assets/Handshake.png";

const tracks = [
  {
    name: "Overall Prize",
    icon: Star,
    description:
      "All projects are considered for the Overall Prize. This track recognizes the most outstanding project at SwampHacks XII, regardless of category, based on innovation, technical achievement, impact, and overall execution.",
  },
  {
    name: "First Timers",
    icon: Handshake,
    description:
      "Everyone starts somewhere. This track celebrates first-time hackers and the projects they build along the way. Teams qualify if at least 50% of their members are attending their first hackathon.",
  },
  {
    name: "Artificial Intelligence & Machine Learning",
    icon: Robot,
    description:
      "Push the boundaries of what's possible. Leverage AI to build projects powered by artificial intelligence, from LLMs and AI agents to computer vision and machine learning that make an impact across any domain.",
  },
  {
    name: "Hardware & Physical Computing",
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
  // {
  //   name: "Health & Wellness",
  //   icon: Apple,
  //   description:
  //     "Support healthier minds and healthier lives. Innovate technology that promotes mental health and physical well-being by supporting emotional wellness and healthy habits. General health-focused projects are also welcome.",
  // },
];

const customModalStyles = {
  overlay: {
    backgroundColor: "rgba(3, 7, 2, 0.62)",
    zIndex: "10000",
  },
  content: {
    top: "50%",
    left: "50%",
    right: "auto",
    bottom: "auto",
    width: "min(calc(100% - 32px), 540px)",
    maxHeight: "82vh",
    marginRight: "-50%",
    padding: "28px",
    overflow: "auto",
    color: "#f4e4cf",
    backgroundColor: "#4d3222",
    border: "3px solid #b98762",
    borderRadius: "4px",
    boxShadow:
      "8px 8px 0 #1d130d, inset 0 0 0 2px #2a1b12",
    outline: "none",
    transform: "translate(-50%, -50%)",
    zIndex: "10001",
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
    <div id="tracks" className="tracks-container">
      <h1 className="tracks-header">Tracks</h1>

      <div className="tracks-body-container">
        <div className="track-container">
          {tracks.map((track, index) => (
            <div
              key={track.name}
              className={`track track-${index}`}
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
            className="modal-close-btn pixel-button pixel-button--danger"
          >
            X
          </button>
          <div>
            <div className="modal-track">
              <img className="modal-track-icon" src={selectedTrack?.icon} />
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
