import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { Lgbutton0Component } from './lgbutton0.component';

describe('Lgbutton0Component', () => {
  let component: Lgbutton0Component;
  let fixture: ComponentFixture<Lgbutton0Component>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ Lgbutton0Component ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(Lgbutton0Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
