ALTER TABLE `wps_ppt_beautify`.`schemes`
CHANGE COLUMN `apply_page` `apply_page` VARCHAR(45) NULL DEFAULT NULL COMMENT '页面类型：\n\ntitle封面页, textt正文页, contents目录页, sectionTitle过渡页, endPage结束页' ;

ALTER TABLE `wps_ppt_beautify`.`schemes_json`
ADD COLUMN `category` VARCHAR(45) NOT NULL AFTER `json`;

ALTER TABLE `wps_ppt_beautify`.`schemes_json`
ADD COLUMN `template_id` VARCHAR(45) NOT NULL AFTER `category`;
