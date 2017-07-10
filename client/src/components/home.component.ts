import {Component, OnInit} from '@angular/core';
import {Jsonp} from '@angular/http';
import {ErrorHandleService} from '../services/error-handle.service';
import 'rxjs/add/operator/map';

@Component({
  selector: 'my-home',
  templateUrl: 'home.component.html',
})
export class HomeComponent implements OnInit {

  stars: number;
  forks: number;
  issues: number;

  constructor(private _jsonp: Jsonp,
              private _eh: ErrorHandleService) {
  }

  ngOnInit() {
    this._jsonp
        .get('https://api.github.com/repos/ronzeidman/ng2-ui-auth?callback=JSONP_CALLBACK')
        .map(res => res.json())
        .subscribe(
            (data) => {
              this.stars = data.data.stargazers_count;
              this.forks = data.data.forks_count;
              this.issues = data.data.open_issues_count;
            },
            (err: any) => this._eh.handleError(err),
        );
  }
}
