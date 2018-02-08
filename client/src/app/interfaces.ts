/**
 * Created by Ron on 02/10/2016.
 */
//copied from server-side, should be in a new shared module and imported by both
export interface IGoogleProfile {
    kind: "plus#personOpenIdConnect";
    gender: string;
    sub: string;
    name: string;
    given_name: string;
    family_name: string;
    profile: string;
    picture: string;
    email: string;
    email_verified: boolean;
    locale: string;
    hd: string;
    error?: Error;
}


export interface ILoginData {
    username: string;
    password: string;
}

export interface ISignupData extends ILoginData {
}

export interface IProfileData {
    displayName: string;
    email: string;
}

export interface IProfileUser {
    username: string;
    email: string;
    displayName: string;
    picture: string;

    l_facebook: boolean;
    l_google: boolean;
    l_linkedin: boolean;
    l_twitter: boolean;
    l_github: boolean;
    l_instagram: boolean;
    l_foursquare: boolean;
    l_yahoo: boolean;
    l_live: boolean;
    l_twitch: boolean;
    l_bitbucket: boolean;
    l_spotify: boolean;
}


export interface IDBUser extends IProfileUser {
    google?: string;
    hash?: string;
}

export interface ISettingsProperty {
  key: string;
  type: string;
  value: any;
  edit: boolean;

  asString(): string;
  asBool(): boolean;
}
