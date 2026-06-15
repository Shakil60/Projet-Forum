USE cinetalk;

INSERT INTO utilisateurs (nom_utilisateur, email, mot_de_passe, sel, role) VALUES
('admin', 'admin@cinetalk.fr', '271b169f31cafd3fb292c0a84b57a6f6b71eb3faba5254e1616db66dbf086c2bf898434a8717ad201ca62aaea5cf818de1493973a4683db6da2d0b48e649dd33', 'a1b2c3d4e5f6a7b8', 'administrateur'),
('alice', 'alice@example.com', '3ad4e699dac4d5f93caa37e7d37c3a34cc23efa06dd7bbbdf9a1bb3441d95e5c8a583a19303fdcc5fb37ce2275adb2195a547efd2976dce9bb890aea4c87f236', '1111111111111111', 'utilisateur'),
('bob', 'bob@example.com', '8f0b8850783aa28c69d4c99e90d572c016fd424a5b36248da3e49720be9192ed6e263c3bcbca402d8da59740d2004d9529fc98045179dc05f5407f30095bab56', '2222222222222222', 'utilisateur'),
('claire', 'claire@example.com', '1bde144784fe47859b31359737dd38903ac5d264e6190eb6be7554fdbda1886248696709411caf95b390e65b6869fb87892b29ff17cc5a4eb46e38bcf58343ed', '3333333333333333', 'utilisateur'),
('david', 'david@example.com', 'a80668a0e1b4ef371ecf3f2699f83c2ed0f6c3a879cb115c7cad7e7d950a4f424e2743aa6518152f2a9c569824411fa3200f156b7d8191b82c15981e7a6822ab', '4444444444444444', 'utilisateur');

INSERT INTO tags (nom) VALUES
('Action'),
('Science-fiction'),
('Drame'),
('Thriller'),
('Comédie'),
('Horreur'),
('Animation'),
('Série'),
('Marvel'),
('Policier');

INSERT INTO fils (titre, contenu, utilisateur_id, etat) VALUES
('La fin d''Inception : la toupie tombe-t-elle ?', 'Après une nouvelle vision du film, je reste persuadé que la dernière scène est volontairement ambiguë. Et vous, qu''en pensez-vous ?', 2, 'ouvert'),
('Quelle est la meilleure série de tous les temps ?', 'Lançons le débat ultime. Pour moi le podium se joue entre Breaking Bad, The Wire et Les Soprano. Donnez vos arguments !', 3, 'ouvert'),
('Le MCU est-il en perte de vitesse ?', 'Depuis Endgame, j''ai l''impression que la qualité des films et séries Marvel est plus inégale. Suis-je le seul à ressentir ça ?', 4, 'ouvert'),
('Films d''horreur à conseiller pour Halloween', 'Je cherche des films d''horreur vraiment efficaces pour une soirée entre amis. Vos meilleures recommandations ?', 5, 'ouvert'),
('Règlement du forum CinéTalk', 'Bienvenue sur CinéTalk ! Merci de rester courtois, d''éviter les spoilers sans avertissement et de respecter les autres membres. Ce fil est fermé aux réponses.', 1, 'ferme'),
('Animation japonaise : vos chefs-d''œuvre', 'Studio Ghibli, Makoto Shinkai, Satoshi Kon... Partagez les films d''animation qui vous ont marqués.', 2, 'ouvert'),
('Débat : Christopher Nolan vs Denis Villeneuve', 'Deux réalisateurs majeurs du cinéma moderne. Lequel a la filmographie la plus solide selon vous ?', 3, 'ouvert'),
('Ancienne discussion archivée', 'Ce sujet n''est plus d''actualité et a été archivé par la modération.', 4, 'archive');

INSERT INTO fil_tags (fil_id, tag_id) VALUES
(1, 2), (1, 4),
(2, 8), (2, 3),
(3, 9), (3, 1),
(4, 6),
(6, 7),
(7, 2), (7, 3),
(8, 3);

INSERT INTO messages (fil_id, utilisateur_id, contenu) VALUES
(1, 3, 'Pour moi elle finit par tomber, on revient bien dans la réalité.'),
(1, 4, 'Nolan a dit que le vrai indice est qu''à la fin Cobb ne regarde même plus la toupie.'),
(1, 5, 'L''ambiguïté volontaire, c''est justement toute la force du film.'),
(2, 2, 'Breaking Bad reste indétrônable, l''évolution de Walter White est parfaite.'),
(2, 5, 'The Wire mérite clairement la première place pour son réalisme social.'),
(2, 1, 'Game of Thrones aurait pu y prétendre... jusqu''à la saison 6.'),
(3, 3, 'Depuis Endgame j''accroche beaucoup moins, l''enjeu émotionnel s''est dilué.'),
(3, 2, 'Trop de séries Disney+ sorties d''un coup, ça a fatigué le public.'),
(4, 4, 'Hereditary m''a vraiment traumatisé, foncez si vous aimez le malaise.'),
(4, 3, 'The Thing de Carpenter, un classique indémodable et toujours aussi efficace.'),
(6, 3, 'Princesse Mononoké, sans la moindre hésitation.'),
(6, 5, 'Your Name est une merveille visuelle et émotionnelle.'),
(7, 4, 'Villeneuve maîtrise l''ambiance comme personne, Dune en est la preuve.'),
(7, 2, 'Nolan reste le roi du concept et du montage, mention spéciale à Interstellar.'),
(5, 1, 'Merci à tous de rester bienveillants et de privilégier le débat constructif.');

INSERT INTO reactions (message_id, utilisateur_id, type) VALUES
(1, 2, 'like'),
(1, 5, 'like'),
(1, 4, 'dislike'),
(3, 3, 'like'),
(3, 2, 'like'),
(3, 4, 'like'),
(4, 5, 'like'),
(4, 4, 'like'),
(4, 3, 'like'),
(4, 1, 'like'),
(5, 2, 'like'),
(5, 3, 'dislike'),
(7, 4, 'like'),
(9, 3, 'like'),
(9, 5, 'like'),
(11, 2, 'like'),
(11, 5, 'like'),
(11, 4, 'like'),
(13, 2, 'like'),
(13, 3, 'like'),
(14, 4, 'like'),
(14, 5, 'dislike');
