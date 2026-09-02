import "./StudentOrgs.css";

import Bat1 from "./assets/bat_without_sign1.png";
import Bat2 from "./assets/bat_without_sign2.png";
import Bat3 from "./assets/bat_without_sign3.png";

const bats = [
  { src: Bat1, top: "78px", left: "3%" },
  { src: Bat2, top: "225px", left: "16%" },
  { src: Bat3, top: "135px", left: "29%" },
  { src: Bat1, top: "185px", left: "70%" },
  { src: Bat2, top: "65px", left: "84%" },
  { src: Bat3, top: "265px", left: "93%" },
];

export default function StudentOrgs() {
  return (
    <div id="studentorgs" className="studentorgs-container">
      <h1 className="studentorgs-header">Student Orgs</h1>
      <p className="studentorgs-coming-soon">Coming soon...</p>

      <div className="studentorgs-background" aria-hidden="true">
        <div className="studentorgs-trees"></div>
      </div>

      <div className="bat-container" aria-hidden="true">
        {bats.map((bat, index) => (
          <img
            className="bat"
            src={bat.src}
            alt=""
            key={index}
            style={{
              top: bat.top,
              left: bat.left,
            }}
          />
        ))}
      </div>
    </div>
  );
}
