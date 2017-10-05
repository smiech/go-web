import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LgFormFieldComponent } from './lg-form-field.component';

describe('LgFormFieldComponent', () => {
  let component: LgFormFieldComponent;
  let fixture: ComponentFixture<LgFormFieldComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LgFormFieldComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LgFormFieldComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
