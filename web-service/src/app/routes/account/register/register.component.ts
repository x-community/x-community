import { Component, OnDestroy } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { _HttpClient } from '@delon/theme';
import { AccountServiceProxy, RegisterRequest } from '@shared/service-proxies/service-proxies';
import { NzMessageService } from 'ng-zorro-antd/message';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'account-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.less'],
})
export class UserRegisterComponent implements OnDestroy {
  constructor(fb: FormBuilder, private router: Router, public accountService: AccountServiceProxy, public msg: NzMessageService) {
    this.form = fb.group({
      mail: [null, [Validators.required, Validators.email]],
      username: [null, [Validators.required, Validators.maxLength(32)]],
      password: [null, [Validators.required, Validators.minLength(8), UserRegisterComponent.checkPassword.bind(this)]],
    });
  }

  // #region fields

  get mail() {
    return this.form.controls.mail;
  }
  get username() {
    return this.form.controls.username;
  }
  get password() {
    return this.form.controls.password;
  }
  form: FormGroup;
  error = '';
  type = 0;
  visible = false;
  status = 'pool';
  progress = 0;
  loading = false;
  passwordProgressMap = {
    ok: 'success',
    pass: 'normal',
    pool: 'exception',
  };

  // #endregion

  // #region get captcha

  count = 0;
  interval$: any;

  static checkPassword(control: FormControl) {
    if (!control) {
      return null;
    }
    const self: any = this;
    self.visible = !!control.value;
    if (control.value && control.value.length > 9) {
      self.status = 'ok';
    } else if (control.value && control.value.length > 5) {
      self.status = 'pass';
    } else {
      self.status = 'pool';
    }

    if (self.visible) {
      self.progress = control.value.length * 10 > 100 ? 100 : control.value.length * 10;
    }
  }

  // #endregion

  submit() {
    this.error = '';
    Object.keys(this.form.controls).forEach((key) => {
      this.form.controls[key].markAsDirty();
      this.form.controls[key].updateValueAndValidity();
    });
    if (this.form.invalid) {
      return;
    }
    this.loading = true;
    const req = new RegisterRequest({ email: this.mail.value, username: this.username.value, password: this.password.value });
    this.accountService.register(req).pipe(finalize(() => { this.loading = false; })).subscribe(_ => {
      this.router.navigateByUrl('/account/register-result?email=' + this.mail.value);
    }, err => {
      this.error = err.error;
    });
  }

  ngOnDestroy(): void {
    if (this.interval$) {
      clearInterval(this.interval$);
    }
  }
}
