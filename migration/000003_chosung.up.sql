CREATE FUNCTION ExtractChosung(input_string varchar(100)) RETURNS varchar(100)
DETERMINISTIC
BEGIN
  DECLARE chosung varchar(100) DEFAULT '';
  DECLARE i int DEFAULT 1;
  
  WHILE i <= LENGTH(input_string) DO
    SET chosung = CONCAT(chosung, 
        CASE
            WHEN SUBSTRING(input_string, i, 1) BETWEEN '가' AND '힣' THEN
                CASE
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '가' AND '깋' THEN 'ㄱ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '까' AND '낗' THEN 'ㄲ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '나' AND '닣' THEN 'ㄴ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '다' AND '딯' THEN 'ㄷ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '따' AND '띻' THEN 'ㄸ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '라' AND '맇' THEN 'ㄹ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '마' AND '밓' THEN 'ㅁ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '바' AND '빟' THEN 'ㅂ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '빠' AND '삫' THEN 'ㅃ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '사' AND '싷' THEN 'ㅅ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '싸' AND '앃' THEN 'ㅆ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '아' AND '잏' THEN 'ㅇ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '자' AND '짛' THEN 'ㅈ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '짜' AND '찧' THEN 'ㅉ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '차' AND '칳' THEN 'ㅊ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '카' AND '킿' THEN 'ㅋ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '타' AND '팋' THEN 'ㅌ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '파' AND '핗' THEN 'ㅍ'
                    WHEN SUBSTRING(input_string, i, 1) BETWEEN '하' AND '힣' THEN 'ㅎ'
                    ELSE ''
                END
            ELSE SUBSTRING(input_string, i, 1)
        END
    );
    SET i = i + 1;
  END WHILE;
  
  RETURN chosung;
END;

CREATE FUNCTION SearchChosung(name varchar(100), input_string varchar(100)) RETURNS tinyint(1)
DETERMINISTIC
BEGIN
  IF input_string REGEXP '^[ㄱㄲㄴㄷㄸㄹㅁㅂㅃㅅㅆㅇㅈㅊㅋㅌㅍㅎ]+$' THEN
    IF ExtractChosung(name) LIKE CONCAT('%', input_string, '%') THEN
      RETURN TRUE;
    ELSE 
      RETURN FALSE;
    END IF;
  ELSE 
    IF name LIKE CONCAT('%', input_string, '%') THEN
      RETURN TRUE;
    ELSE
      RETURN FALSE;
    END IF;
  END IF;
END;