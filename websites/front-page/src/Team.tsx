import { memo, useState } from "react";
import "./Team.css";

import Abby from "./assets/headshots/Abigail_Moore_cropped.jpg";
import Chris from "./assets/headshots/Christian_Cardenas_cropped.jpg";
import Gabby from "./assets/headshots/Gabby_Houser_cropped.jpg";
import Lu from "./assets/headshots/Lu_Ighodalo_cropped.jpg";
import Hector from "./assets/headshots/Hector_Morel_cropped.jpg";
import Jared from "./assets/headshots/Jared_Cohen_cropped.png";
import Hieu from "./assets/headshots/Hieu_Nguyen_cropped.png";
import Jessie from "./assets/headshots/Jessie_Lin_cropped.jpg";
import Shivani from "./assets/headshots/Shivani_Chandrasekar_cropped.jpg";
import Lucia from "./assets/headshots/Lucia_Novo_cropped.jpg";
import Micah from "./assets/headshots/Micah_Tam_cropped.jpg";
import Mason from "./assets/headshots/Mason_Alexander_cropped.jpg";
import Kyle from "./assets/headshots/Kyle_Scarmack_cropped.jpg";
import Steve from "./assets/headshots/Steve_Sajeev_cropped.jpg";
import Oliver from "./assets/headshots/Oliver_Jen_cropped.png";
import Lara from "./assets/headshots/Lara_Afont_cropped.jpg";
import Nam from "./assets/headshots/Nam_Tran_cropped.png";
import Maria from "./assets/headshots/Maria_Davis_cropped.jpg";
import Christopher from "./assets/headshots/Christopher_Silva_cropped.jpg";
import Matilde from "./assets/headshots/Matilde_Gillia_cropped.jpg";
import Kian from "./assets/headshots/Kian_Hidalgo_cropped.jpg";
import Belinda from "./assets/headshots/Belinda_M_Lopez_cropped.jpg";
import Masha from "./assets/headshots/Masha_N_Belyaeva_cropped.jpg";
import Leah from "./assets/headshots/Leah_Zabad_cropped.jpg";
import Dhivya from "./assets/headshots/Dhivya_Kumar_cropped.png";
import Leandra from "./assets/headshots/Leandra_Fleitas_cropped.png";
import Ayushi from "./assets/headshots/Ayushi_Srivastava_cropped.jpg";
import Kaiden from "./assets/headshots/Kaiden_Joy_cropped.jpg";
import Ivette from "./assets/headshots/Ivette_S_Hernandez_cropped.jpg";
import Will from "./assets/headshots/William_Chi_cropped.jpg";

import bat1 from "./assets/bat_without_sign1.png";
import bat2 from "./assets/bat_without_sign2.png";
import bat3 from "./assets/bat_without_sign3.png";

const teamMembers = [
  {
    name: "Christian Cardenas",
    linkedin: "https://www.linkedin.com/in/christian-cardenas-5a07302aa/",
    role: "Organizers",
    img: Chris,
  },
  {
    name: "Gabby Houser",
    linkedin: "https://www.linkedin.com/in/gabriela-houser/",
    role: "Organizers",
    img: Gabby,
  },
  {
    name: "Lu Ighodalo",
    linkedin: "https://www.linkedin.com/in/luighodalo/",
    role: "Logistics & Operations",
    isExec: true,
    img: Lu,
  },
  {
    name: "Hector Morel",
    linkedin: "https://www.linkedin.com/in/hector-morel/",
    role: "Logistics & Operations",
    isExec: true,
    img: Hector,
  },
  {
    name: "Jared Cohen",
    linkedin: "https://www.linkedin.com/in/jared-cohen-86533b397",
    role: "Finance",
    isExec: true,
    img: Jared,
  },
  {
    name: "Abby Moore",
    linkedin: "https://www.linkedin.com/in/abigail-g-moore",
    role: "Technology",
    isExec: true,
    img: Abby,
  },
  {
    name: "Hieu Nguyen",
    linkedin: "https://www.linkedin.com/in/hieutnguyendev/",
    role: "Technology",
    isExec: true,
    img: Hieu,
  },
  {
    name: "Jessie Lin",
    linkedin: "https://www.linkedin.com/in/jessielin04/",
    role: "Marketing",
    isExec: true,
    img: Jessie,
  },
  {
    name: "Shivani Chandrasekar",
    linkedin: "https://www.linkedin.com/in/shivani-chandrasekar-901357311/",
    role: "Marketing",
    isExec: true,
    img: Shivani,
  },
  {
    name: "Lucia Novo",
    linkedin: "https://www.linkedin.com/in/lucia-novo-6aa862221/",
    role: "Logistics & Operations",
    isExec: false,
    img: Lucia,
  },
  {
    name: "Micah Tam",
    linkedin: "https://www.linkedin.com/in/micahtam/",
    role: "Finance",
    isExec: false,
    img: Micah,
  },
  {
    name: "Mason Alexander",
    linkedin: "https://www.linkedin.com/in/mason-alexander-",
    role: "Finance",
    isExec: false,
    img: Mason,
  },
  {
    name: "Kyle Scarmack",
    linkedin: "https://www.linkedin.com/in/kylescarmack/",
    role: "Logistics & Operations",
    isExec: false,
    img: Kyle,
  },
  {
    name: "Steve Sajeev",
    linkedin: "https://www.linkedin.com/in/stevesajeev/",
    role: "Technology",
    isExec: false,
    img: Steve,
  },
  {
    name: "Oliver Jen",
    linkedin: "https://www.linkedin.com/in/oliver-jen-ufl/",
    role: "Logistics & Operations",
    isExec: false,
    img: Oliver,
  },
  {
    name: "Lara Afont",
    linkedin: "https://www.linkedin.com/in/laraafont/",
    role: "Technology",
    isExec: false,
    img: Lara,
  },
  {
    name: "Nam Tran",
    linkedin: "https://www.linkedin.com/in/nam2k5/",
    role: "Technology",
    isExec: false,
    img: Nam,
  },
  {
    name: "Maria Davis",
    linkedin: "https://www.linkedin.com/in/maria-davis-/",
    role: "Finance",
    isExec: false,
    img: Maria,
  },
  {
    name: "Christopher Silva",
    linkedin: "https://www.linkedin.com/in/christopher-silva-4b9902325/",
    role: "Technology",
    isExec: false,
    img: Christopher,
  },
  {
    name: "Matilde Gillia",
    linkedin: "https://www.linkedin.com/in/matilde-gillia/",
    role: "Logistics & Operations",
    isExec: false,
    img: Matilde,
  },
  {
    name: "Kian Hidalgo",
    linkedin: "https://www.linkedin.com/in/kian-hidalgo-141b96338/",
    role: "Logistics & Operations",
    isExec: false,
    img: Kian,
  },
  {
    name: "Belinda Morales Lopez",
    linkedin: "https://www.linkedin.com/in/belindamoraleslopez",
    role: "Marketing",
    isExec: false,
    img: Belinda,
  },
  {
    name: "Masha Nicole Belyaeva",
    linkedin: "https://www.linkedin.com/in/masha-belyaeva-513b97311",
    role: "Logistics & Operations",
    isExec: false,
    img: Masha,
  },
  {
    name: "Leah Zabad",
    linkedin: "https://www.linkedin.com/in/leah-zabad",
    role: "Technology",
    isExec: false,
    img: Leah,
  },
  {
    name: "Dhivya Kumar",
    linkedin: "https://www.linkedin.com/in/dhivyadkumar",
    role: "Logistics & Operations",
    isExec: false,
    img: Dhivya,
  },
  {
    name: "Leandra Fleitas",
    linkedin: "https://www.linkedin.com/in/leandra-f-709a14348/",
    role: "Marketing",
    isExec: false,
    img: Leandra,
  },
  {
    name: "Ayushi Srivastava",
    linkedin: "https://www.linkedin.com/in/ayusrivastava/",
    role: "Marketing",
    isExec: false,
    img: Ayushi,
  },
  {
    name: "Kaiden Joy",
    linkedin: "https://www.linkedin.com/in/kaiden-joy-934a84261/",
    role: "Finance",
    isExec: false,
    img: Kaiden,
  },
  {
    name: "Ivette Saldana Hernandez",
    linkedin: "https://www.linkedin.com/in/ivettesaldanahernandez/",
    role: "Technology",
    isExec: false,
    img: Ivette,
  },
  {
    name: "William Chi",
    linkedin: "https://www.linkedin.com/in/chiwilliam/",
    role: "Technology",
    isExec: false,
    img: Will,
  },
];

const groupedMembers = [
  {
    role: "Organizers",
    members: teamMembers.filter((member) => member.role === "Organizers"),
  },
  {
    role: "Logistics & Operations",
    members: teamMembers.filter(
      (member) => member.role === "Logistics & Operations",
    ),
  },
  {
    role: "Finance",
    members: teamMembers.filter((member) => member.role === "Finance"),
  },
  {
    role: "Technology",
    members: teamMembers.filter((member) => member.role === "Technology"),
  },
  {
    role: "Marketing",
    members: teamMembers.filter((member) => member.role === "Marketing"),
  },
];

type TeamMember = (typeof teamMembers)[number];

const TeamMemberCard = memo(function TeamMemberCard({
  member,
  compact,
}: {
  member: TeamMember;
  compact: boolean;
}) {
  return (
    <a className="team-member-card" href={member.linkedin} target="_blank">
      <img
        className={`team-member-image ${compact ? "is-compact" : ""}`.trim()}
        src={member.img}
        alt={member.name}
      />
      <p className="team-member-name">
        {member.name}
        {member.isExec && " (Executive)"}
      </p>
    </a>
  );
});

export default function Team() {
  const [activeRoleIndex, setActiveRoleIndex] = useState(0);

  const handleTabSelect = (index: number) => {
    setActiveRoleIndex(index);
  };

  const currentGroup = groupedMembers[activeRoleIndex];

  return (
    <div className="team-container">
      <h1 className="team-header">Meet the XII Team!</h1>

      <div className="background">
        <div className="layer vegetation"></div>
        <div className="layer canopy"></div>
        <div className="layer trees-back1"></div>
        <div className="layer foreground1"></div>

        <img
          className="bat"
          src={bat1}
          style={{
            left: "10%",
          }}
        />
        <img
          className="bat"
          src={bat2}
          style={{
            top: "15%",
            right: "10%",
          }}
        />
        <img
          className="bat"
          src={bat3}
          style={{
            top: "-12%",
            right: "15%",
          }}
        />
        <img
          className="bat"
          src={bat2}
          style={{
            top: "-12%",
            left: "30%",
          }}
        />

        <img
          className="bat"
          src={bat3}
          style={{
            left: "15%",
            bottom: "80%",
          }}
        />
      </div>

      <div className="team-area">
        <div className="team-content" aria-label="Team members by role">
          <div className="team-tabs" role="tablist" aria-label="Team roles">
            {groupedMembers.map((group, index) => {
              const isActive = index === activeRoleIndex;

              return (
                <button
                  key={group.role}
                  type="button"
                  role="tab"
                  aria-selected={isActive}
                  aria-controls={`team-panel-${index}`}
                  className={`team-tab ${isActive ? "is-active" : ""}`.trim()}
                  onClick={() => handleTabSelect(index)}
                >
                  {group.role}
                </button>
              );
            })}
          </div>

          <div className="team-role-panels">
            {groupedMembers.map((group, index) => {
              const isActive = index === activeRoleIndex;

              return (
                <div
                  className={`team-role-group ${isActive ? "is-active" : ""}`}
                  key={group.role}
                  role="tabpanel"
                  aria-hidden={!isActive}
                >
                  <h2 className="team-role-title">{group.role}</h2>

                  <div className="team-role-members">
                    {group.members.map((member) => (
                      <TeamMemberCard
                        key={`${group.role}-${member.name}`}
                        member={member}
                        compact={group.members.length <= 2}
                      />
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
