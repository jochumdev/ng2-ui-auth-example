import {Injectable, ViewContainerRef} from '@angular/core';
import {ToastsManager} from 'ng2-toastr/ng2-toastr';
import {Response} from '@angular/http';
/**
 * Created by ronze on 10/9/2016.
 */
@Injectable()
export class ErrorHandleService {
  private _savedError: any;
  private _savedSeverity: string;
  private _savedMessage: string;

  constructor(private _toastr: ToastsManager) {
  }

  setRootViewContainerRef(vcr: ViewContainerRef) {
    this._toastr.setRootViewContainerRef(vcr);
  }

  saveMessage(severity: string, message: string) {
    this._savedSeverity = severity;
    this._savedMessage = message;
  }

  showSavedMessage() {
    if (this._savedSeverity == "" || this._savedMessage == "") {
      return
    }

    if (this._savedSeverity == "error") {
      this._toastr.error(this._savedMessage);
    } else if (this._savedSeverity == "warning") {
      this._toastr.warning(this._savedMessage);
    } else if (this._savedSeverity == "info") {
      this._toastr.info(this._savedMessage);
    } else if (this._savedSeverity == "success") {
      this._toastr.success(this._savedMessage);
    }

    this._savedSeverity = "";
    this._savedMessage = "";
  }

  saveError(err: any) {
    this._savedError = err;
  }

  showSavedError() {
    if (this._savedError == null) {
      return
    }

    this.handleError(this._savedError);
    this._savedError = null;
  }

  handleError(err: any) {
      if (typeof err === 'string') {
          this._toastr.error(err)
      } else if (err instanceof Response) {
          const res: Response = err;
          if (res.text() && res.text() !== res.statusText) {
              this._toastr.error(res.text(), res.statusText);
          } else {
              this._toastr.error(res.statusText);
          }
      } else if (err && err.message) {
          this._toastr.error(err.message);
      } else if (err) {
          this._toastr.error(err.toString());
      } else {
          this._toastr.error('An unknown error has occurred');
      }
  }
}
