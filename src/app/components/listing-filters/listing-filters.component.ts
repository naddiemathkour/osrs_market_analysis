import { Component, Output } from '@angular/core';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatSliderModule } from '@angular/material/slider';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { IFilters } from '../../interfaces/filters.interface';

@Component({
  selector: 'app-listing-filters',
  standalone: true,
  imports: [
    MatButtonToggleModule,
    MatSliderModule,
    MatSelectModule,
    MatCheckboxModule,
  ],
  templateUrl: './listing-filters.component.html',
  styleUrl: './listing-filters.component.scss',
})
export class ListingFiltersComponent {
  @Output() filterValues: IFilters = {} as IFilters;
}
