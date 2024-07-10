import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AlchingViewComponent } from './alching-view.component';

describe('AlchingViewComponent', () => {
  let component: AlchingViewComponent;
  let fixture: ComponentFixture<AlchingViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AlchingViewComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AlchingViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
