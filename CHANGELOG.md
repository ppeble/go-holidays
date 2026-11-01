# Go Holidays CHANGELOG

## 1.2.0

* Update to [v9.3.0 definitions](https://github.com/holidays/definitions/releases/tag/v9.3.0), in parity with the Ruby `holidays` gem [11.8.0](https://rubygems.org/gems/holidays/versions/11.8.0).
* Support the optional `excluded_regions` key on holiday rules, so a holiday listed for a parent region can skip specific subregions. [#34](https://github.com/holidays/go-holidays/pull/34)
* Exclude the 17 states without a paid Juneteenth state holiday (including `us_hi`) from the formal `us` Juneteenth and list it as informal there instead. [holidays issue-444](https://github.com/holidays/holidays/issues/444)
* Fix `au_tas` getting New Year's Day on a weekend 1 January and a duplicate on the following Monday. [definitions issue-37](https://github.com/holidays/definitions/issues/37)
* Pin `tr` Ramazan and Kurban Bayramı to the Diyanet dates for 2031-2035 and add the second Ramazan Bayramı of 2033. [#33](https://github.com/holidays/go-holidays/pull/33)
* Add `cl` Día Nacional de los Pueblos Indígenas. [#27](https://github.com/holidays/go-holidays/pull/27)
* Extend the lunar tables to 2100 and return `ErrUnsupportedYear` past them. [#31](https://github.com/holidays/go-holidays/pull/31)
* Fix the `vn` lunar table to use Vietnam time instead of the Korean table, and the `kr` 2042 row so month 11 starts on 2042-12-12. [#29](https://github.com/holidays/go-holidays/pull/29), [#30](https://github.com/holidays/go-holidays/pull/30)

## 1.1.0

* Update to [v9.2.0 definitions](https://github.com/holidays/definitions/releases/tag/v9.2.0), in parity with the Ruby `holidays` gem [11.7.0](https://rubygems.org/gems/holidays/versions/11.7.0).
* Add `ma` Eid al-Fitr and Eid al-Adha, calculated from the Hijri calendar.
* Add `my` Hari Raya Puasa and Hari Raya Haji, calculated from the Hijri calendar and observed on the following Monday when they fall on a Sunday.

## 1.0.0

* Initial release of the Go port, in parity with the Ruby `holidays` gem [11.6.0](https://rubygems.org/gems/holidays/versions/11.6.0), shipping [v9.1.0 definitions](https://github.com/holidays/definitions/releases/tag/v9.1.0).
