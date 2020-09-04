import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AccountServiceProxy, VerifyAccountRequest } from '@shared/service-proxies/service-proxies';
import { NzMessageService } from 'ng-zorro-antd/message';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'account-active',
  templateUrl: './active.component.html',
})
export class UserActiveComponent implements OnInit {

  loading = false;
  code: string;
  error: string;

  constructor(route: ActivatedRoute, private router: Router, private accountService: AccountServiceProxy) {
    this.code = route.snapshot.queryParams.code;
  }

  ngOnInit(): void {
    if (!this.code || !this.code.length) {
      this.router.navigateByUrl('/');
    }
    this.loading = true;
    this.accountService.verify(new VerifyAccountRequest({ code: this.code })).pipe(finalize(() => { this.loading = false; })).subscribe(resp => {
      console.log(resp);
      this.router.navigateByUrl('/');
    }, err => {
      this.error = err.error;
    });
  }
}
