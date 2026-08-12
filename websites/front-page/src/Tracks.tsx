import "./Tracks.css";
import pencil from "./assets/Pencil.png";
import controller from "./assets/Controller.png";
import star from "./assets/Star.png";
import leaf from "./assets/Leaf.png";
import handshake from "./assets/Handshake.png";
import Gator from "./assets/Gator.gif";
import { useState } from "react";
import Modal from "react-modal";

const tracks = [
  {
    name: "General Track (Overall Prize)",
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
