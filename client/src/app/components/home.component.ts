import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ErrorHandleService } from '../services/error-handle.service';
import 'rxjs/add/operator/map';

@Component({
  selector: 'my-home',
  templateUrl: 'home.component.html',
})
export class HomeComponent implements OnInit {

  stars: number;
  forks: number;
  issues: number;

  constructor(private _http: HttpClient,
    private _eh: ErrorHandleService) {
  }

  ngOnInit() {}

  // ngOnInit() {
  //   this._http.jsonp('https://api.github.com/repos/ronzeidman/ng2-ui-auth', 'JSONP_CALLBACK')
  //     .subscribe(
  //     (data) => {
  //       this.stars = data.stargazers_count;
  //       this.forks = data.forks_count;
  //       this.issues = data.open_issues_count;
  //     },
  //     (err: any) => this._eh.handleError(err),
  //   );
  // }
}
