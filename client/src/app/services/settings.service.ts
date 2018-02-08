import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { ConfigService } from 'ng2-ui-auth';
import { Observable } from 'rxjs/Observable';
import { ErrorHandleService } from './error-handle.service';
import { ISettingsProperty } from '../interfaces';

import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';

interface IOauth2Options {
  clientId?: string;
}

/**
 * Created by pcdummy on 07/03/2017.
 */
export const SettingAuthAllowSignup = "app_auth.AllowSignup";
export const SettingAuthGoogleClientID = "app_auth.GoogleClientID";
export const SettingAuthFacebookClientID = "app_auth.FacebookClientID";
export const SettingAuthGithubClientID = "app_auth.GithubClientID";
export const SettingTwitterEnabled = "app_auth.TwitterEnabled";

export class SettingsProperty implements ISettingsProperty {
  key: string;
  type: string;
  value: any;
  edit: boolean;

  constructor(key: string, type: string, value: any, edit: boolean) {
    this.key = key;
    this.type = type;
    this.value = value;
    this.edit = edit;
  }

  asString(): string {
    return <string>this.value;
  }

  asBool(): boolean {
    return <boolean>this.value;
  }
}

@Injectable()
export class SettingsService {
  settings: Map<string, SettingsProperty>;

  constructor(private _http: HttpClient,
    private _eh: ErrorHandleService,
    private _ac: ConfigService) {
    this.settings = new Map<string, SettingsProperty>();
  }

  fetchSettings() {
    this._http.get('/api/settings/v1/')
      .map((response: any) => <Array<SettingsProperty>>response)
      .catch((err: HttpErrorResponse) => {
        this._eh.handleError(err);
        return Observable.throw('Server error');
      })
      .subscribe((settings: Array<SettingsProperty>) => {
        settings.forEach((setting: SettingsProperty) => {
          this.settings.set(setting.key, new SettingsProperty(
            setting.key, setting.type, setting.value, setting.edit
          ));
        });

        if (this.settings.has(SettingAuthGoogleClientID)) {
          var options: IOauth2Options = <IOauth2Options>this._ac.options.providers['google'];
          options.clientId = this.settings.get(SettingAuthGoogleClientID).asString();
        }
        if (this.settings.has(SettingAuthFacebookClientID)) {
          var options: IOauth2Options = <IOauth2Options>this._ac.options.providers['facebook'];
          options.clientId = this.settings.get(SettingAuthFacebookClientID).asString();
        }
        if (this.settings.has(SettingAuthGithubClientID)) {
          var options: IOauth2Options = <IOauth2Options>this._ac.options.providers['github'];
          options.clientId = this.settings.get(SettingAuthGithubClientID).asString();
        }
      });
  }

  hasProvider(provider: string): boolean {
    if (provider == 'google') {
      if (this.settings.has(SettingAuthGoogleClientID)) {
        return this.settings.get(SettingAuthGoogleClientID).asString() != "";
      }
    } else if (provider == "facebook") {
      if (this.settings.has(SettingAuthFacebookClientID)) {
        return this.settings.get(SettingAuthFacebookClientID).asString() != "";
      }
    } else if (provider == "github") {
      if (this.settings.has(SettingAuthGithubClientID)) {
        return this.settings.get(SettingAuthGithubClientID).asString() != "";
      }
    } else if (provider == "twitter") {
      if (this.settings.has(SettingTwitterEnabled)) {
        return this.settings.get(SettingTwitterEnabled).asBool();
      }
    }

    return false;
  }
}
