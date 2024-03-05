import { Server, Member, Profile } from './models';

export type ServerWithMembersWithProfiles = Server & {
    members: (Member & { profile: Profile })[];
};
