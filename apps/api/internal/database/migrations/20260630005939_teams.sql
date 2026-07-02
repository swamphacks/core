-- +goose Up
alter table teams drop column hackathon_id;
alter table teams drop constraint teams_owner_id_fkey;
alter table teams add constraint teams_owner_id_fkey
foreign key (owner_id) references users(id) on delete cascade, alter column owner_id set not null;
alter table teams add constraint unique_owner_id unique (owner_id);

alter table team_invitations add constraint unique_team_id unique (team_id);
alter table team_invitations drop column invited_email;
alter table team_invitations rename column invited_by_user_id to inviter_id;
alter table team_invitations drop column invited_user_id;
alter table team_invitations drop column status;

-- +goose Down
alter table teams add hackathon_id text not null references hackathons(id);
alter table teams drop constraint teams_owner_id_fkey;
alter table teams add constraint teams_owner_id_fkey
foreign key (owner_id) references users(id) on delete set null, alter column owner_id drop not null;
alter table teams drop constraint unique_owner_id;

alter table team_invitations drop constraint unique_team_id;
alter table team_invitations add invited_email text not null;
alter table team_invitations add invited_user_id uuid references users on delete cascade;
alter table team_invitations rename column inviter_id to invited_by_user_id;
alter table team_invitations add column status team_invitation_status default 'pending'::team_invitation_status not null;