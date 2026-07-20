import "./StudentOrgs.css";
import UFSIT from "./assets/clubs/UFSIT.svg";
import CSU from "./assets/clubs/CSU.svg";
import CSK from "./assets/clubs/CSK.svg";

export default function StudentOrgs() {
  return (
    <div className="studentorgs-container">
      <h1 className="studentorgs-header">Student Orgs</h1>

      <div className="background">
        <div className="layer vegetation"></div>
        <div className="layer layer-canopy"></div>
        <div className="layer trees-back1"></div>
        <div className="layer trees-back2"></div>
        <div className="layer foreground1"></div>
      </div>

      <div className="bat-container">
        <img
          className="bat"
          src={UFSIT}
          style={{
            top: "100px",
            left: "10%",
          }}
        />
        <img
          className="bat"
          src={CSU}
          style={{
            top: "200px",
            left: "20%",
          }}
        />
        <img
          className="bat"
          src={CSK}
          style={{
            top: "100px",
            left: "50%",
          }}
        />
      </div>
    </div>
  );
}
