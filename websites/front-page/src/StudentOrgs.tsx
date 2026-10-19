import "./Sponsors.css";
import "./StudentOrgs.css";

import UFSIT from "./assets/clubs/UFSIT.svg";
import GUD from "./assets/clubs/GUD.svg";
import FES from "./assets/clubs/FES.svg";
import WECE from "./assets/clubs/WECE.svg";
import DSI from "./assets/clubs/DSI.svg";
import ColorStack from "./assets/clubs/ColorStack.svg";
import BADAS from "./assets/clubs/BADAS.svg";
import GatorVR from "./assets/clubs/GatorVR.svg";
import SPCB from "./assets/clubs/SPCB.svg";
import SEC from "./assets/clubs/SEC.svg";
import GatorAI from "./assets/clubs/GatorAI.svg";
import OSC from "./assets/clubs/OSC.svg";

import Bat1 from "./assets/bat_without_sign1.png";
import Bat2 from "./assets/bat_without_sign2.png";
import Bat3 from "./assets/bat_without_sign3.png";

const decorativeBats = [Bat1, Bat2, Bat3];

function BatColumn() {
  return (
    <div className="bat-container" aria-hidden="true">
      {decorativeBats.map((src) => (
        <img className="bat" src={src} alt="" key={src} />
      ))}
    </div>
  );
}

type StudentOrgPartner = {
  name: string;
  logo: string;
  url?: string;
};

// student orgs
const studentOrgPartners: StudentOrgPartner[] = [
  { name: "WECE", logo: WECE, url: "https://www.instagram.com/wece_uf/" },
  { name: "DSI", logo: DSI, url: "https://www.instagram.com/uf_dsi/" },
  { name: "UFSIT", logo: UFSIT, url: "https://www.instagram.com/uf.sit/" },
  { name: "GUD", logo: GUD, url: "https://www.instagram.com/gatoruserdesign/"},
  { name: "FES", logo: FES, url: "https://www.instagram.com/floridaengineeringsociety/"},
  { name: "ColorStack", logo: ColorStack, url: "https://www.instagram.com/colorstackuf/" },
  { name: "BADAS", logo: BADAS, url: "https://www.instagram.com/badassocietyclub/" },
  { name: "GatorVR", logo: GatorVR, url: "https://www.instagram.com/ufgatorvr/" },
  { name: "SPCB", logo: SPCB, url: "https://www.instagram.com/pcbuildinguf/" },
  { name: "SEC", logo: SEC, url: "https://www.instagram.com/ufsec/" },
  { name: "GatorAI", logo: GatorAI, url: "https://www.instagram.com/uf_gatorai/" },
  { name: "OSC", logo: OSC, url: "https://www.instagram.com/uf_osc/" },
];

export default function StudentOrgs() {
  return (
    <div id="studentorgs" className="studentorgs-container">
      <h1 className="sponsors-header studentorgs-header">Student Orgs</h1>
      <div className="studentorgs-layout">
        <BatColumn />
        <div className="sponsor-tiers studentorgs-partners" aria-label="Student organization partners">
          <div className="sponsor-grid sponsor-grid--medium">
            {studentOrgPartners.map((partner) =>
              partner.url ? (
                <a
                  className="sponsor-card sponsor-card--medium studentorgs-partner"
                  href={partner.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={partner.name}
                  key={partner.name}
                >
                  <img src={partner.logo} alt={partner.name} />
                </a>
              ) : (
                <div className="sponsor-card sponsor-card--medium studentorgs-partner" key={partner.name}>
                  <img src={partner.logo} alt={partner.name} />
                </div>
              ),
            )}
          </div>
        </div>

        <BatColumn />
      </div>

      <div className="studentorgs-background" aria-hidden="true">
        <div className="studentorgs-trees"></div>
      </div>
    </div>
  );
}
