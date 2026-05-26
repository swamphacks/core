import Gabby from "./assets/pfps/Gabby.png";
import Chris from "./assets/pfps/Chris.png";

const team = [
  {
    position: "Organizers",
    members: [
      { name: "Christian Cardenas", pic: Gabby, linkedin: "" },
      { name: "Gabby Houser", pic: Chris, linkedin: "" },
    ],
  },
  {
    position: "Marketing",
    members: [
      { name: "Christian Cardenas", pic: Gabby, linkedin: "" },
      { name: "Gabby Houser", pic: Chris, linkedin: "" },
    ],
  },
  {
    position: "Finance",
    members: [
      { name: "Christian Cardenas", pic: Gabby, linkedin: "" },
      { name: "Gabby Houser", pic: Chris, linkedin: "" },
    ],
  },
  {
    position: "Logistics",
    members: [
      { name: "Christian Cardenas", pic: Gabby, linkedin: "" },
      { name: "Gabby Houser", pic: Chris, linkedin: "" },
    ],
  },
  {
    position: "Technology",
    members: [
      { name: "Abby Moore (Executive)", pic: Chris, linkedin: "" },
      {
        name: "Hieu Nguyen (Executive)",
        pic: Gabby,
        linkedin: "https://www.linkedin.com/in/hieutnguyendev",
      },
    ],
  },
];

export default function Team() {
  return (
    <div>
      <p className="tab-title">SH XII Team</p>
      <div className="team-grid team-container">
        {team.map((group) => (
          <section className="team-group" key={group.position}>
            <h2 className="team-position">{group.position}</h2>
            <div className="team-members">
              {group.members.map((member) => (
                <article className="team-card" key={member.name}>
                  <div className="avatar">
                    <img
                      src={member.pic}
                      className="profile-pic"
                      alt={member.name}
                    />
                  </div>
                  <p className="team-name">{member.name}</p>
                  <a
                    className="profile-linkedin"
                    target="_blank"
                    href={member.linkedin}
                  >
                    <i className="nes-icon linkedin profile-linkedin"></i>
                  </a>
                </article>
              ))}
            </div>
          </section>
        ))}
      </div>
    </div>
  );
}
